package grpc

// const bufSize = 1024 * 1024
//
// var lis *bufconn.Listener
//
// func init() {
// 	lis = bufconn.Listen(bufSize)
// 	s := grpc.NewServer()
// 	userService.RegisterUserServiceServer(s, &server{})
// 	go func() {
// 		if err := s.Serve(lis); err != nil {
// 			log.Fatalf("Server exited with error: %v", err)
// 		}
// 	}()
// }
//
// func bufDialer(context.Context, string) (net.Conn, error) {
// 	return lis.Dial()
// }
//
// func TestUserServer_SignUp(t *testing.T) {
// 	type fields struct {
// 		UnimplementedUserServiceServer userService.UnimplementedUserServiceServer
// 		log                            telemetry.AppLogger
// 		cfg                            *config.App
// 		svc                            delivery.User
// 		jwtManager                     auth.JWTManager
// 	}
// 	type args struct {
// 		ctx     context.Context
// 		request *pbUser.SignUpRequest
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    *emptypb.Empty
// 		wantErr bool
// 	}{
// 		{},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			u := &UserServer{
// 				UnimplementedUserServiceServer: tt.fields.UnimplementedUserServiceServer,
// 				log:                            tt.fields.log,
// 				cfg:                            tt.fields.cfg,
// 				svc:                            tt.fields.svc,
// 				jwtManager:                     tt.fields.jwtManager,
// 			}
// 			got, err := u.SignUp(tt.args.ctx, tt.args.request)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("SignUp() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("SignUp() got = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
