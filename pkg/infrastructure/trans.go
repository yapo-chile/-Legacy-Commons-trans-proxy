package infrastructure

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"

	"github.com/eapache/go-resiliency/retrier"
	"gitlab.com/yapo_team/legacy/commons/trans-proxy/pkg/domain"
	"gitlab.com/yapo_team/legacy/commons/trans-proxy/pkg/interfaces/loggers"
	"gitlab.com/yapo_team/legacy/commons/trans-proxy/pkg/interfaces/repository/services"
)

// trans-proxy struct definition
type trans-proxy struct {
	conf            TransConf
	logger          loggers.Logger
	allowedCommands []string
}

// textProtocolTransFactory is a auxiliar struct to create trans-proxy on demand
type textProtocolTransFactory struct {
	conf            TransConf
	logger          loggers.Logger
	allowedCommands []string
}

// NewTextProtocolTransFactory initialize a services.TransFactory
func NewTextProtocolTransFactory(
	conf TransConf,
	logger loggers.Logger,
) services.TransFactory {
	return &textProtocolTransFactory{
		conf:            conf,
		logger:          logger,
		allowedCommands: strings.Split(conf.AllowedCommands, "|"),
	}
}

// MakeTransHandler initialize a services.TransHandler on demand
func (t *textProtocolTransFactory) MakeTransHandler() services.TransHandler {
	return &trans-proxy{
		conf:            t.conf,
		logger:          t.logger,
		allowedCommands: t.allowedCommands,
	}
}

// SendCommand use a socket connection to send commands to trans-proxy port
func (handler *trans-proxy) SendCommand(cmd string, trans-proxyParams []domain.TransParams) (map[string]string, error) {
	respMap := make(map[string]string)
	// check if the command is allowed; if not, return error
	valid := handler.isAllowedCommand(cmd)
	if !valid {
		err := fmt.Errorf(
			"invalid command - commands allowed: %s",
			handler.allowedCommands,
		)
		respMap["error"] = err.Error()
		handler.logger.Error(err.Error())
		return respMap, err
	}
	conn, err := handler.connect()
	if err != nil {
		handler.logger.Error("Error connecting to trans-proxy: %s\n", err.Error())
		return respMap, fmt.Errorf("Error connecting with trans-proxy server")
	}
	defer conn.Close() //nolint: errcheck, megacheck

	// initiate the context so the request can timeout
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(handler.conf.Timeout)*time.Second,
	)
	defer cancel()

	respMap, err = handler.sendWithContext(ctx, conn, cmd, trans-proxyParams)
	if err != nil {
		handler.logger.Error("Error Sending command %s: %s\n", cmd, err)
	}

	return respMap, err
}

// isAllowedCommand checks if the given command can be sent to trans-proxy
func (handler *trans-proxy) isAllowedCommand(cmd string) bool {
	for _, allowedCommand := range handler.allowedCommands {
		if allowedCommand == cmd {
			return true
		}
	}
	return false
}

// connect returns a connection to the trans-proxy client.
// Retries to connect after retryAfter time if the connection times out
func (handler *trans-proxy) connect() (net.Conn, error) {
	// initiate the retrier that will handle retry reconnect if the connection dies
	r := retrier.New(
		[]time.Duration{
			time.Duration(handler.conf.RetryAfter) * time.Second},
		nil,
	)
	var conn net.Conn
	// set the function that starts the connection
	err := r.Run(func() error {
		var e error
		conn, e = net.DialTimeout(
			"tcp",
			fmt.Sprintf(
				"%s:%d",
				handler.conf.Host,
				handler.conf.Port,
			),
			time.Duration(handler.conf.Timeout)*time.Second,
		)
		return e
	})
	return conn, err
}

// sendWithContext sends the message to trans-proxy but is cancelable via a context.
// The context timeout specified how long the caller can wait
// for the trans-proxy to respond
func (handler *trans-proxy) sendWithContext(
	ctx context.Context,
	conn io.ReadWriteCloser,
	cmd string,
	args []domain.TransParams,
) (map[string]string, error) {
	var resp map[string]string
	errChan := make(chan error, 1)

	// starts the go routine that sends the message and retrieves the response and error, if any.
	// it communicates any error through errChan
	go func() {
		errChan <- func() error {
			var err error
			resp, err = handler.send(conn, cmd, args)
			return err
		}()
	}()

	select {
	case <-ctx.Done():
		// closing the connection here interrupts the send function, in the gorouting above, if it
		// is waiting on reading from or writing to the connection.
		err := conn.Close()
		if err != nil {
			handler.logger.Error("Error Closing connection to trans-proxy after ctx done: %s\n", err.Error())
		}
		// wait for the goroutine to return and ignore the error
		<-errChan
		// return the context error: the operation timed out.
		return nil, ctx.Err()
	case err := <-errChan:
		// in this case the send function returned before
		// the timeout of the context.
		return resp, err
	}
}

func (handler *trans-proxy) send(conn io.ReadWriter, cmd string, args []domain.TransParams) (map[string]string, error) {
	// Check greeting.
	reader := bufio.NewReader(conn)
	line, err := reader.ReadSlice('\n')
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(line, []byte("220 Welcome.\n")) {
		return nil, fmt.Errorf("trans-proxy: unexpected greeting: %q", line)
	}

	buf := make([]byte, 0)
	// Send command to Trans.
	buf = appendCmd(buf, cmd, args)
	if _, err = conn.Write(buf); err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	_, err = buffer.ReadFrom(reader)
	if err != nil {
		return nil, err
	}

	buf, encodingErr := charmap.ISO8859_1.NewDecoder().Bytes(bytes.TrimSuffix(buffer.Bytes(), []byte("end\n")))
	if encodingErr != nil {
		handler.logger.Debug("Latin 1 expected, encoding error: %s\n", encodingErr.Error())
	}
	respMap, err := TransResponse(buf).Map()
	if err != nil {
		return respMap, fmt.Errorf("error parsing response: %s", err.Error())
	}
	return respMap, nil
}

// appendCmd Appends the command to the buffer. For the command format, see:
// https://scmcoord.com/wiki/Trans#Protocol
func appendCmd(buf []byte, cmd string, args []domain.TransParams) []byte {
	buf = append(buf, "cmd:"...)
	buf = append(buf, cmd...)
	buf = append(buf, '\n')
	for _, param := range args {
		key := param.Key
		if value, ok := param.Value.(string); ok {
			if param.Blob {
				if decoded, err := base64.StdEncoding.DecodeString(value); err == nil {
					value = string(decoded)
					buf = append(buf, "blob:"...)
					buf = strconv.AppendInt(buf, int64(len(value)), 10)
					buf = append(buf, ':')
					buf = append(buf, key...)
					buf = append(buf, '\n')
					buf = append(buf, value...)
					buf = append(buf, '\n')
				}
				continue
			}
			key, err := charmap.ISO8859_1.NewEncoder().String(key)
			if err != nil {
				continue
			}
			value, err = charmap.ISO8859_1.NewEncoder().String(value)
			if err != nil {
				continue
			}
			buf = append(buf, key...)
			buf = append(buf, ':')
			buf = append(buf, value...)
			buf = append(buf, '\n')
		}
	}
	buf = append(buf, "commit:1"...)
	buf = append(buf, "\nend\n"...)
	return buf
}
