package client

var DefaultClientMap = map[string]Client{}

func NewClients(client Client) {
	DefaultClientMap[client.Name] = client
}
