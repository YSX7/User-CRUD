package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	service "vuegolang/proto"
)

type server struct {
	service.UnimplementedAddServiceServer
	db *pgx.Conn
}

func main() {
	if errDB != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	lis, err := net.Listen("tcp", ":4040")
	if err != nil {
		panic(err)
	}
	srv := grpc.NewServer()
	service.RegisterAddServiceServer(srv, &server{db: conn})
	reflection.Register(srv)

	if e := srv.Serve(lis); e != nil {
		panic(err)
	}
}

func (s server) Auth(ctx context.Context, request *service.Request) (*emptypb.Empty, error) {
	var (
		login, pass = request.GetLogin(), request.GetPassword()
		empty       = new(emptypb.Empty)
		dbPass      string
	)

	fmt.Println("request")
	if len(login) == 0 || len(pass) == 0 {
		return empty, status.Error(codes.InvalidArgument, "не хватает аргументов, один/оба пустые")
	}

	// Формируем запрос к базе
	query, _, errQ := goqu.From("users").Select("password").Where(goqu.Ex{"username": login}).ToSQL()
	if errQ != nil {
		panic(errQ)
	}

	// Запрашиваем
	err := s.db.QueryRow(context.Background(), query).Scan(&dbPass)

	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return empty, status.Error(codes.InvalidArgument, "неверный логин/пароль")
		default:
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(1)
		}
	}

	errHash := bcrypt.CompareHashAndPassword([]byte(dbPass), []byte(pass))
	if errHash != nil {
		return empty, status.Error(codes.InvalidArgument, "неверный логин/пароль")
	}

	return empty, nil
}
