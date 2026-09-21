package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	flowpb "github.com/cilium/cilium/api/v1/flow"
	observerpb "github.com/cilium/cilium/api/v1/observer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func main() {
	address := flag.String("address", "127.0.0.1:4245", "Hubble Relay address")
	namespace := flag.String("namespace", "", "client-side source or destination namespace filter")
	limit := flag.Uint64("limit", 20, "maximum flow responses")
	timeout := flag.Duration("timeout", 10*time.Second, "request timeout")
	caFile := flag.String("ca", "", "CA certificate file for Relay TLS")
	certFile := flag.String("cert", "", "client certificate file for Relay mTLS")
	keyFile := flag.String("key", "", "client private key file for Relay mTLS")
	insecure := flag.Bool("insecure", false, "use plaintext gRPC; only for an explicitly plaintext local endpoint")
	flag.Parse()

	if flag.NArg() != 1 {
		fatal("usage: hubble-probe [flags] status|nodes|namespaces|flows")
	}

	conn, err := dial(*address, *caFile, *certFile, *keyFile, *insecure)
	if err != nil {
		fatal(err.Error())
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	client := observerpb.NewObserverClient(conn)

	switch flag.Arg(0) {
	case "status":
		response, err := client.ServerStatus(ctx, &observerpb.ServerStatusRequest{})
		printProto(response, err)
	case "nodes":
		response, err := client.GetNodes(ctx, &observerpb.GetNodesRequest{})
		printProto(response, err)
	case "namespaces":
		response, err := client.GetNamespaces(ctx, &observerpb.GetNamespacesRequest{})
		printProto(response, err)
	case "flows":
		stream, err := client.GetFlows(ctx, &observerpb.GetFlowsRequest{Number: *limit})
		if err != nil {
			fatal(err.Error())
		}
		seen := uint64(0)
		for seen < *limit {
			response, receiveErr := stream.Recv()
			if receiveErr != nil {
				if receiveErr == io.EOF {
					break
				}
				fatal(receiveErr.Error())
			}
			seen++
			if *namespace != "" && !matchesNamespace(response.GetFlow(), *namespace) {
				continue
			}
			printProto(response, nil)
		}
	default:
		fatal("unknown command: " + flag.Arg(0))
	}
}

func dial(address, caFile, certFile, keyFile string, insecure bool) (*grpc.ClientConn, error) {
	if insecure {
		return grpc.NewClient(address, grpc.WithInsecure())
	}
	if caFile == "" {
		return nil, fmt.Errorf("TLS is enabled by default; provide -ca, or use -insecure only for a plaintext local endpoint")
	}
	caData, err := os.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM(caData) {
		return nil, fmt.Errorf("could not parse CA certificate: %s", caFile)
	}
	config := &tls.Config{RootCAs: roots, MinVersion: tls.VersionTLS12}
	if certFile != "" || keyFile != "" {
		if certFile == "" || keyFile == "" {
			return nil, fmt.Errorf("-cert and -key must be supplied together")
		}
		certificate, certErr := tls.LoadX509KeyPair(certFile, keyFile)
		if certErr != nil {
			return nil, certErr
		}
		config.Certificates = []tls.Certificate{certificate}
	}
	return grpc.NewClient(address, grpc.WithTransportCredentials(credentials.NewTLS(config)))
}

func matchesNamespace(flow *flowpb.Flow, namespace string) bool {
	if flow == nil {
		return false
	}
	return flow.GetSource().GetNamespace() == namespace || flow.GetDestination().GetNamespace() == namespace
}

func printProto(value proto.Message, err error) {
	if err != nil {
		fatal(err.Error())
	}
	data, marshalErr := protojson.MarshalOptions{UseProtoNames: true}.Marshal(value)
	if marshalErr != nil {
		fatal(marshalErr.Error())
	}
	fmt.Println(string(data))
}

func fatal(message string) {
	fmt.Fprintln(os.Stderr, "hubble-probe:", strings.TrimSpace(message))
	os.Exit(1)
}
