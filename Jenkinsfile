node {
    def root = tool name: 'Go 1.8', type: 'go'
    withEnv(["GOROOT=${root}", "PATH+GO=${root}/bin"]) {
       sh 'go version'

       ws("${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/") {
           withEnv(["GOPATH=${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"]) {
               env.PATH="${GOPATH}/bin:$PATH"

               def app

               stage('Clone repository') {
                   /* Let's make sure we have the repository cloned to our workspace */
                   checkout scm
               }

               stage('Pre Test'){
                    sh 'go version'
               }

               stage('Build'){
                   sh """cd $GOPATH/src/github.com/cjburchell/reefstatus-go && go build -o main ."""
               }

               stage('Build image') {
                   /* This builds the actual image; synonymous to
                    * docker build on the command line */

                   app = docker.build("cjburchell/reefstatus")
               }
           }
       }
    }
}