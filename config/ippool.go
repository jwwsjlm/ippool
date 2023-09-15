package config

import "sync"

type IpPool struct {
	Hash   sync.Map
	IpChan chan string
}

func NewMap(i int) *IpPool {
	ip := IpPool{
		Hash:   sync.Map{},
		IpChan: make(chan string, i),
	}

	return &ip

}
func (i *IpPool) AddIP(ip string) {
	i.IpChan <- ip
}
func (i *IpPool) GetIP() string {
	p := <-i.IpChan
	return p
}
