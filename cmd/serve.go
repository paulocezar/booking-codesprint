// Copyright © 2016 Paulo Cezar <paulocezar.ufg@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"path"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/paulocezar/booking-codesprint/passions"
	"github.com/paulocezar/booking-codesprint/search"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Launches a rpc service on localhost:13337 and search webserver on http://localhost:1337",
	Run:   serve,
}

var (
	rpcPort      int
	port         int
	swaggerDir   string
	databasePath string
)

func init() {
	RootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntVarP(&rpcPort, "rpc-port", "r", 13337, "port where the rpc server will be available")
	serveCmd.Flags().IntVarP(&port, "port", "p", 1337, "port where the server will be available")
	serveCmd.Flags().StringVarP(&swaggerDir, "swagger-dir", "s", "passions", "path to the directory which contains swagger definitions")
	serveCmd.Flags().StringVarP(&databasePath, "db-path", "d", "datasets/destinations.csv", "path to the csv file with the data for our Simple Search Engine")
}

func serve(_ *cobra.Command, _ []string) {
	err := startRPCServer()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := http.NewServeMux()
	mux.HandleFunc("/swagger/", serveSwagger)

	gw, err := newGateway(ctx)
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/", gw)

	log.Println("passions-hacked REST API available on", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), mux)
}

func startRPCServer() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%v", rpcPort))
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	srv, err := search.NewSimpleSearchServer(databasePath)
	if err != nil {
		return err
	}

	passions.RegisterPassionServicesServer(s, srv)

	go s.Serve(l)
	log.Println("passions-hacked RPC server started on", rpcPort)
	return nil
}

func newGateway(ctx context.Context, opts ...runtime.ServeMuxOption) (http.Handler, error) {
	mux := runtime.NewServeMux(opts...)
	dialOpts := []grpc.DialOption{grpc.WithInsecure()}
	err := passions.RegisterPassionServicesHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%v", rpcPort), dialOpts)
	if err != nil {
		return nil, err
	}
	return mux, nil
}

func serveSwagger(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.URL.Path, ".swagger.json") {
		log.Printf("Not Found: %s", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	log.Printf("Serving %s", r.URL.Path)
	p := strings.TrimPrefix(r.URL.Path, "/swagger/")
	p = path.Join(swaggerDir, p)
	http.ServeFile(w, r, p)
}
