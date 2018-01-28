node {
    stage('Clone repository') {
        /* Let's make sure we have the repository cloned to our workspace */
        checkout scm
    }

    String goPath = "/go/src/github.com/cjburchell/reefstatus-go"
    String workspacePath =  "/volume1/Storage/jenkins-data/workspace/ReefStatus"

    stage('Build'){
        docker.image('golang:1.8.0-alpine').inside("-v ${workspacePath}:${goPath}"){
         sh """cd ${goPath} && go build -o main ."""
        }
    }

    stage('Build image') {
        docker.build("cjburchell/reefstatus")
    }

     stage('Push image') {
            /* Finally, we'll push the image with two tags:
             * First, the incremental build number from Jenkins
             * Second, the 'latest' tag.
             * Pushing multiple tags is cheap, as all the layers are reused. */
            docker.withRegistry('https://registry.hub.docker.com', 'docker-hub-credentials') {
                app.push("${env.BUILD_NUMBER}")
                app.push("latest")
            }
     }
}