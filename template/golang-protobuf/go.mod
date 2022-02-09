module handler

go 1.16

replace handler/function => ./function

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/openfaas/templates-sdk v0.0.0-20200723110415-a699ec277c12
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd
	golang.org/x/sys v0.0.0-20220207234003-57398862261d // indirect
	google.golang.org/genproto v0.0.0-20220207185906-7721543eae58 // indirect
	google.golang.org/grpc v1.44.0
	google.golang.org/protobuf v1.27.1 // indirect
)
