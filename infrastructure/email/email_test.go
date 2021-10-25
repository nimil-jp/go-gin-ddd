package email

import "testing"

func Test_email_Send(t *testing.T) {
	type args struct {
		to   string
		body Body
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "default",
			args: args{
				to: "test@test.com",
				body: DefaultBody{
					Title: "テスト",
					Body:  "メール送信テスト",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := email{}
			if err := e.Send(tt.args.to, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
