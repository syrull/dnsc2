package pkg

import (
	"time"
)

type Client struct {
	Id          int
	RemoteIp    string
	MachineId   string
	LastUpdated int64
	CreatedAt   int64
}

func NewClient(id int, remoteIp string, lastUpdated int, createdAt int) *Client {
	now := time.Now()
	return &Client{
		Id:          id,
		RemoteIp:    remoteIp,
		LastUpdated: now.Unix(),
		CreatedAt:   now.Unix(),
	}
}

// func (c *Client) insert(db *sql.DB) error {
// 	insertStatement := fmt.Sprintf(`
// 		INSERT INTO client(remote_ip, last_updated, created_at)
// 		VALUES('%s', '%d', '%d');
// 	`, c.RemoteIp, c.LastUpdated, c.CreatedAt)
// 	if _, err := db.Exec(insertStatement); err != nil {
// 		return err
// 	}
// 	return nil
// }
