IMAGE_NAME=rickyslash/fiber-ecommerce:1.0.0

docker build -t ${IMAGE_NAME} .

docker push ${IMAGE_NAME}