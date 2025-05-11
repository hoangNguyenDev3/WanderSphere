#!/bin/bash

set -e

# Go to the project root directory
cd "$(dirname "$0")/.."

# Clean up existing wrong directory structure
rm -rf pkg/types/proto/pb
mkdir -p pkg/types/proto/pb

# Re-generate protos
export PATH=$PATH:$HOME/go/bin
echo "Generating Authpost proto files..."
mkdir -p pkg/types/proto/pb/authpost
protoc --proto_path=. \
  --go_out=pkg/types/proto/pb/authpost \
  --go_opt=paths=import \
  --go-grpc_out=pkg/types/proto/pb/authpost \
  --go-grpc_opt=paths=import \
  pkg/types/proto/authpost.proto

echo "Generating Newsfeed proto files..."
mkdir -p pkg/types/proto/pb/newsfeed
protoc --proto_path=. \
  --go_out=pkg/types/proto/pb/newsfeed \
  --go_opt=paths=import \
  --go-grpc_out=pkg/types/proto/pb/newsfeed \
  --go-grpc_opt=paths=import \
  pkg/types/proto/newsfeed.proto

echo "Generating Newsfeed Publishing proto files..."
mkdir -p pkg/types/proto/pb/newsfeed_publishing
protoc --proto_path=. \
  --go_out=pkg/types/proto/pb/newsfeed_publishing \
  --go_opt=paths=import \
  --go-grpc_out=pkg/types/proto/pb/newsfeed_publishing \
  --go-grpc_opt=paths=import \
  pkg/types/proto/newsfeed_publishing.proto

# Apply the compatibility patch
for grpc_file in $(find pkg/types/proto/pb -name "*_grpc.pb.go"); do
  # Remove StaticMethod() calls and SupportPackageIsVersion9 references
  sed -i 's/append(\[\]grpc.CallOption{grpc.StaticMethod()}, opts...)/opts/g' "$grpc_file"
  sed -i 's/const _ = grpc.SupportPackageIsVersion9/\/\/ Using compatible gRPC version/g' "$grpc_file"
  echo "Patched $grpc_file for compatibility"
done

echo "All proto files generated and fixed successfully!" 