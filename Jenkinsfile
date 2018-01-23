pipeline {
        agent {
            docker {
                image 'golang'
                args '-p 3000:3000'
                customWorkspace '${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}'
            }
        }

        environment {
                GOPATH = '${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}'
                PATH =  "${GOPATH}/bin:$PATH"
        }


        stages {
                stage('Checkout'){
                    echo 'Checking out SCM'
                    checkout scm
                }

                stage('Pre Test'){
                    echo 'Pulling Dependencies'
                    sh 'go version'
                }

                stage('Test'){

                    //List all our project files with 'go list ./... | grep -v /vendor/ | grep -v github.com | grep -v golang.org'
                    //Push our project files relative to ./src
                    sh 'cd $GOPATH/src/github.com/cjburchell/reefstatus-go/ && go list ./... | grep -v /vendor/ > projectPaths'

                    //Print them with 'awk '$0="./src/"$0' projectPaths' in order to get full relative path to $GOPATH
                    def paths = sh returnStdout: true, script: """awk '\$0="./src/"\$0' projectPaths"""

                    echo 'Vetting'

                    sh """cd $GOPATH && go tool vet ${paths}"""

                    echo 'Linting'
                    sh """cd $GOPATH && golint ${paths}"""

                    echo 'Testing'
                    sh """cd $GOPATH && go test -race -cover ${paths}"""
                }

                stage('Build') {
                    echo 'Building Executable'

                    sh """cd $GOPATH/src/github.com/cjburchell/reefstatus-go/ && go build -o service"""
                }
        }
}