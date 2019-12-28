package ram

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/git-sim/tc/app/domain/repo"
)

func TestWithStringStruct(t *testing.T) {

	type MsgStruct struct {
		id         uint64
		sender     string
		recipients []string
		subject    string
		body       []byte
	}

	gr := NewGenericRepo()
	count, _ := gr.RetrieveCount()
	if count != 0 {
		t.Error("Count not initialized to 0")
	}

	msgs := make([]MsgStruct, 3)
	for i, m := range msgs {
		m.id = uint64(i)
		m.sender = "abc@mail.com"
		m.recipients = []string{"def@mail.com", "ghi@mail.com"}
		m.subject = fmt.Sprintf("test subject %d", i)
		m.body = ([]byte)(fmt.Sprintf("test message %d", i))
		msgs[i] = m
		gr.Create(repo.GenericKeyT(m.id), m)
	}

	count, _ = gr.RetrieveCount()
	if len(msgs) != count {
		t.Errorf("count expected %d got %d", len(msgs), count)
	}

	for i, m := range msgs {

		val, err := gr.Retrieve(repo.GenericKeyT(i))
		if err != nil {
			t.Errorf("retreive failed for id %d, err %s", i, err)
		}
		msgStruct, ok := val.(MsgStruct)
		if ok {
			if !reflect.DeepEqual(msgStruct, m) {
				t.Errorf("structs don't match idx %d", i)
			}
		} else {
			t.Errorf("failed to convert back to msgStruct")
		}
	}

	for i, _ := range msgs {
		gr.Delete(repo.GenericKeyT(i))
	}

	count, _ = gr.RetrieveCount()
	if count != 0 {
		t.Error("Count not 0 after delete")
	}
}
