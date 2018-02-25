package server

import (
	"context"
	"log"

	"github.com/nrocco/ide/pkg/ide"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) RefreshCtags(ctx context.Context, in *RefreshCtagsRequest) (*RefreshCtagsReply, error) {
	project, err := ide.LoadProject(in.Directory)
	if err != nil {
		return &RefreshCtagsReply{}, err
	}

	go func() {
		log.Printf("Start generating ctags file: %s", project.CtagsFile())
		if err := project.RefreshCtags(); err != nil {
			log.Printf("Error generating ctags file: %s", err)
		} else {
			log.Printf("Done generating ctags file: %s", project.CtagsFile())
		}
	}()

	return &RefreshCtagsReply{File: project.CtagsFile()}, nil
}

// NewServer returns a new grpc server with registered Ide services
func NewServer() *grpc.Server {
	s := grpc.NewServer()
	RegisterServerServer(s, &server{})

	return s
}
