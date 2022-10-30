/**
 * @Author: fuxiao
 * @Email: 576101059@qq.com
 * @Date: 2021/8/26 1:59 下午
 * @Desc: TODO
 */

package http

import (
	"os"
	"strings"

	"github.com/inkrtech/tencent-im-http/internal"
)

var contentTypeToFileSuffix = map[string]string{
	"application/x-001":              ".001",
	"text/h323":                      ".323",
	"drawing/907":                    ".907",
	"audio/x-mei-aac":                ".acp",
	"audio/aiff":                     ".aif",
	"text/asa":                       ".asa",
	"text/asp":                       ".asp",
	"audio/basic":                    ".au",
	"application/vnd.adobe.workflow": ".awf",
	"application/x-bmp":              ".bmp",
	"application/x-c4t":              ".c4t",
	"application/x-cals":             ".cal",
	"application/x-netcdf":           ".cdf",
	"application/x-cel":              ".cel",
	"application/x-g4":               ".cg4",
	"application/x-cit":              ".cit",
	"text/xml":                       ".cml",
	"application/x-cmx":              ".cmx",
	"application/pkix-crl":           ".crl",
	"application/x-csi":              ".csi",
	"application/x-cut":              ".cut",
	"application/x-dbm":              ".dbm",
}

type Download struct {
	request *Request
}

func NewDownload(c *Client) *Download {
	return &Download{
		request: NewRequest(c),
	}
}

// Download download a file from the network address to the local.
func (d *Download) Download(url, dir string, filename ...string) (string, error) {
	resp, err := d.request.request(MethodGet, url)
	if err != nil {
		return "", err
	}

	var path string

	if len(filename) > 0 {
		path = strings.TrimRight(dir, string(os.PathSeparator)) + string(os.PathSeparator) + filename[0]
	} else {
		path = d.genFilePath(resp, dir)
	}

	if err = internal.SaveToFile(path, resp.ReadBytes()); err != nil {
		return "", err
	}

	return path, nil
}

// genFilePath generate file path based on response content type
func (d *Download) genFilePath(resp *Response, dir string) string {
	path := strings.TrimRight(dir, string(os.PathSeparator)) + string(os.PathSeparator) + internal.RandStr(16)

	if suffix := internal.GetFileType(resp.ReadBytes()); suffix != "" {
		path += "." + suffix
	}

	if internal.Exists(path) {
		return d.genFilePath(resp, dir)
	}

	return path
}
