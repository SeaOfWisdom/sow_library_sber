package rest

import (
	"fmt"

	"github.com/SeaOfWisdom/sow_library/src/service/storage"
)

type WorkReq struct {
	Work *storage.Work `json:"work"`
}

func (r *WorkReq) Validate() error {
	if r.Work == nil {
		return fmt.Errorf("work is null")
	}
	if r.Work.Name == "" { // TODO
		return fmt.Errorf("wrong work name: %s", r.Work.Name)
	}
	if r.Work.Annotation == "" { // TODO
		return fmt.Errorf("wrong work annotation: %s", r.Work.Annotation)
	}
	if r.Work.Content == nil { // TODO
		return fmt.Errorf("work content is null")
	}
	if r.Work.Content.WorkData == "" { // TODO
		return fmt.Errorf("work content data is null")
	}
	return nil
}

type WorkResp struct {
	Status storage.WorkStatus `json:"work_status" example:"WORK_UNDER_PRE_REVIEW"`
}

type PublishWorkDataResp struct {
	Tags []string `json:"tags" eaxmple:"[подводный спорт, моноласт]"`
}
