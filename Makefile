# push docker image to aws ecr
push-docker-image:
	@docker tag zolaris-backend-app-stage 864981729345.dkr.ecr.ap-south-1.amazonaws.com/zolaris-go-app:latest
	@docker push 864981729345.dkr.ecr.ap-south-1.amazonaws.com/zolaris-go-app:latest
