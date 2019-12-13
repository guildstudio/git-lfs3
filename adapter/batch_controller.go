package adapter

import (
	"encoding/json"
	"time"

	"github.com/ikmski/git-lfs3/usecase"
)

const (
	contentMediaType = "application/vnd.git-lfs"
	metaMediaType    = contentMediaType + "+json"
)

// BatchRequest is ...
type BatchRequest struct {
	Operation string           `json:"operation"`
	Transfers []string         `json:"transfers,omitempty"`
	Ref       Ref              `json:"ref,omitempty"`
	Objects   []*ObjectRequest `json:"objects"`
}

// ObjectRequest is ...
type ObjectRequest struct {
	Oid      string `json:"oid"`
	Size     int64  `json:"size"`
	User     string
	Password string
	Repo     string
}

// Ref is ...
type Ref struct {
	Name string `json:"name"`
}

// BatchResponse is ...
type BatchResponse struct {
	Transfer string            `json:"transfer,omitempty"`
	Objects  []*ResponseObject `json:"objects"`
}

// ResponseObject is ...
type ResponseObject struct {
	Oid     string           `json:"oid"`
	Size    int64            `json:"size"`
	Actions map[string]*Link `json:"actions"`
	Error   *ObjectError     `json:"error,omitempty"`
}

// Link is ...
type Link struct {
	Href      string            `json:"href"`
	Header    map[string]string `json:"header,omitempty"`
	ExpiresAt time.Time         `json:"expires_at,omitempty"`
}

// ObjectError is ...
type ObjectError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// BatchController is ...
type BatchController interface {
	Batch(ctx Context)
}

type batchController struct {
	BatchService usecase.BatchService
}

// NewBatchController is ...
func NewBatchController(s usecase.BatchService) BatchController {
	return &batchController{
		BatchService: s,
	}
}

func (c *batchController) Batch(ctx Context) {

	req := parseBatchRequest(ctx)

	result, err := c.BatchService.Batch(req)

	if err != nil {

	}

	res := convertBatchResponse(result)

	json, err := json.Marshal(res)
	if err != nil {

	}

	ctx.Header("Content-Type", metaMediaType)
	ctx.JSON(200, json)
}

func parseBatchRequest(ctx Context) *usecase.BatchRequest {

	var br BatchRequest

	data, err := ctx.GetRawData()
	if err != nil {
		return convertBatchRequest(&br)
	}

	err = json.Unmarshal(data, &br)
	if err != nil {
		return convertBatchRequest(&br)
	}

	for i := 0; i < len(br.Objects); i++ {
		br.Objects[i].User = ctx.Param("user")
		br.Objects[i].Repo = ctx.Param("repo")
	}

		return convertBatchRequest(&br)
}

func convertBatchResponse(result *usecase.BatchResult) *BatchResponse {

	var objs []*ResponseObject

	res := &BatchResponse{
		Transfer: "basic",
		Objects:  objs,
	}

	for _, batchObj := range result.Objects {

		header := make(map[string]string)
		header["Accept"] = contentMediaType

		obj := &ResponseObject{}

		if batchObj.MetaExists {

		}

		if batchObj.ObjectExists {
			obj.Actions["download"] = &Link{
				Href:   "https://hoge",
				Header: header,
			}
		} else {
			obj.Actions["upload"] = &Link{
				Href:   "https://hoge",
				Header: header,
			}
		}

		objs = append(objs, obj)
	}

	return res
}

func convertBatchRequest(req *BatchRequest) *usecase.BatchRequest {

    var objs []*usecase.ObjectRequest

    for _, o := range req.Objects {

        item := &usecase.ObjectRequest {
            Oid: o.Oid,
            Size: o.Size,
            User: o.User,
            Password: o.Password,
            Repo: o.Repo,
        }

        objs = append(objs, item)
    }

    br := &usecase.BatchRequest {
        Operation: req.Operation,
        Transfers: req.Transfers,
        Ref: req.Ref.Name,
        Objects: objs,
    }

    return br
}