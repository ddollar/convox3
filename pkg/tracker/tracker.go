package tracker

import (
	"fmt"
	"os"
	"strings"

	"github.com/convox/console/api/model"
	"github.com/convox/console/pkg/settings"
	"github.com/convox/console/pkg/storage"
	raven "github.com/convox/raven-go"
	"github.com/convox/stdapi"
)

var (
	Model = model.New(storage.New("dynamo"))
)

func init() {
	raven.SetDSN("https://478e8f9c3d19413792b81c374592b1f0:0d2f2f0e19c6488988e5151a1bc33ef4@sentry.io/1299449")
	raven.SetEnvironment(settings.Host)
	raven.SetRelease(os.Getenv("RELEASE"))
}

type causer interface {
	Cause() error
}

func CaptureApiError(err error, c *stdapi.Context) string {
	fmt.Printf("ERROR: %s\n", err)

	is := []raven.Interface{}

	is = append(is, raven.NewHttp(c.Request()))

	if uid, ok := c.Get("uid").(string); ok {
		if u, _ := Model.UserGet(uid); u != nil {
			is = append(is, &raven.User{
				ID:    u.ID,
				Email: u.Email,
			})
		}
	}

	return captureError(err, is)
}

func CaptureError(err error) string {
	fmt.Printf("ERROR: %s\n", err)

	is := []raven.Interface{}

	return captureError(err, is)
}

func captureError(err error, is []raven.Interface) string {
	packet := raven.NewPacket(err.Error(), append(is, raven.NewException(err, raven.GetOrNewStacktrace(err, 1, 3, raven.IncludePaths())))...)

	esp := strings.Split(err.Error(), ": ")
	packet.Fingerprint = []string{"{{default}}", strings.TrimSpace(esp[len(esp)-1])}

	tags := map[string]string{
		"version": settings.Version,
	}

	id, _ := raven.Capture(packet, tags)

	return id
}
