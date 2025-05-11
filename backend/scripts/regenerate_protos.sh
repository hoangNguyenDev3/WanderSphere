#!/bin/bash

set -e

# Go to the project root directory
cd "$(dirname "$0")/.."

# Ensure we use a compatible gRPC version
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.0
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0

# Clean existing generated files
rm -rf pkg/types/proto/pb
mkdir -p pkg/types/proto/pb

# Make sure we don't have old .pb.go files in the proto directory
rm -f pkg/types/proto/*.pb.go

# Find all proto files and their corresponding directories
AUTHPOST_PROTO="pkg/types/proto/authpost.proto"
NEWSFEED_PROTO="pkg/types/proto/newsfeed.proto"
NEWSFEED_PUBLISHING_PROTO="pkg/types/proto/newsfeed_publishing.proto"

echo "Generating code for Authpost proto"
mkdir -p pkg/types/proto/pb/authpost
protoc --proto_path=. \
  --go_out=pkg/types/proto/pb/authpost \
  --go_opt=paths=import \
  --go-grpc_out=pkg/types/proto/pb/authpost \
  --go-grpc_opt=paths=import \
  ${AUTHPOST_PROTO}

echo "Generating code for Newsfeed proto"
mkdir -p pkg/types/proto/pb/newsfeed
protoc --proto_path=. \
  --go_out=pkg/types/proto/pb/newsfeed \
  --go_opt=paths=import \
  --go-grpc_out=pkg/types/proto/pb/newsfeed \
  --go-grpc_opt=paths=import \
  ${NEWSFEED_PROTO}

echo "Generating code for Newsfeed Publishing proto"
mkdir -p pkg/types/proto/pb/newsfeed_publishing
protoc --proto_path=. \
  --go_out=pkg/types/proto/pb/newsfeed_publishing \
  --go_opt=paths=import \
  --go-grpc_out=pkg/types/proto/pb/newsfeed_publishing \
  --go-grpc_opt=paths=import \
  ${NEWSFEED_PUBLISHING_PROTO}

# Move the files to their correct locations
echo "Moving files to correct locations..."
for service in authpost newsfeed newsfeed_publishing; do
  # Copy the files from the deep path to the top-level directory
  cp -r pkg/types/proto/pb/$service/github.com/hoangNguyenDev3/WanderSphere/backend/pkg/types/proto/pb/$service/* pkg/types/proto/pb/$service/
  # Remove the deep directory structure
  rm -rf pkg/types/proto/pb/$service/github.com
done

# Apply a patch for compatibility with grpc v1.45.0
for grpc_file in $(find pkg/types/proto/pb -name "*_grpc.pb.go"); do
  # Remove StaticMethod() calls and SupportPackageIsVersion9 references
  sed -i 's/append(\[\]grpc.CallOption{grpc.StaticMethod()}, opts...)/opts/g' "$grpc_file"
  sed -i 's/const _ = grpc.SupportPackageIsVersion9/\/\/ Using compatible gRPC version/g' "$grpc_file"
  echo "Patched $grpc_file for compatibility"
done

echo "All proto files generated successfully!" 