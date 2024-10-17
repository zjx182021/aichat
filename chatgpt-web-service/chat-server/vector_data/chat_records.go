package vector_data

import (
	"chatgpt-web-service/pkg/config"
	"context"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
)

const CHAT_RECORDS = "chat_records"

type ChatRecord struct {
	ID  string
	KVs map[string]string
}
type IChatRecordsData interface {
	UpsertData(ctx context.Context, list []*ChatRecord) error
	QueryData(ctx context.Context, text map[string][]string) (id string, score float32, err error)
}

type chatRecordsData struct {
	config   *config.Config
	vectorDB *tcvectordb.Client
}

func NewChatRecordsData(config *config.Config, vectorDB *tcvectordb.Client) IChatRecordsData {
	return &chatRecordsData{
		config:   config,
		vectorDB: vectorDB,
	}
}
func (data *chatRecordsData) UpsertData(ctx context.Context, list []*ChatRecord) error {
	database := data.config.VectorDB.Database
	collection := CHAT_RECORDS
	coll := data.vectorDB.Database(database).Collection(collection)
	documentList := make([]tcvectordb.Document, 0, len(list))
	for _, l := range list {
		doc := tcvectordb.Document{
			Id: l.ID,
		}
		doc.Fields = make(map[string]tcvectordb.Field, len(l.KVs))
		for k, v := range l.KVs {
			doc.Fields[k] = tcvectordb.Field{Val: v}
		}
		documentList = append(documentList, doc)
	}
	_, err := coll.Upsert(ctx, documentList)
	if err != nil {
		return err
	}
	return nil
}
func (data *chatRecordsData) QueryData(ctx context.Context, text map[string][]string) (id string, score float32, err error) {
	database := data.config.VectorDB.Database
	collection := CHAT_RECORDS
	coll := data.vectorDB.Database(database).Collection(collection)
	result, err := coll.SearchByText(ctx, text, &tcvectordb.SearchDocumentParams{
		Params: &tcvectordb.SearchDocParams{Ef: 100},
		Limit:  1,
	})
	if err != nil {
		return "", 0, err
	}
	if len(result.Documents) > 0 && len(result.Documents[0]) > 0 {
		doc := result.Documents[0][0]
		return doc.Id, doc.Score, nil

	}
	return "", 0, nil
}
