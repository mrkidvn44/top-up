package mapper

import authpb "top-up-api/proto/auth"

func ToAuthRequest(token string, userID uint64) *authpb.AuthenticateServiceRequest {
	return &authpb.AuthenticateServiceRequest{
		TokenString: token,
		UserId:      userID,
	}
}
