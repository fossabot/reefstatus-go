node {
    stage('Clone repository') {
        /* Let's make sure we have the repository cloned to our workspace */
        checkout scm
    }

    stage('Build'){
        docker.image('golang').inside("-v ${env.WORKSPACE}:/go/src/github.com/cjburchell/reefstatus-go "){

         // Debugging
         sh 'echo GOPATH: $GOPATH'
         sh "ls -al /go/src/github.com/cjburchell/reefstatus-go"
         sh "cd /go/src/github.com/cjburchell/reefstatus-go"
         sh "pwd"

         sh """cd /go/src/github.com/cjburchell/reefstatus-go && go build -o main ."""
        }
    }

    stage('Build image') {
        /* This builds the actual image; synonymous to
         * docker build on the command line */

        docker.build("cjburchell/reefstatus")
    }
}