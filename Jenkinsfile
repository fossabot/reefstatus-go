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
                    steps {
                        echo 'Checking out SCM'
                        checkout scm
                    }
                }

                stage('Pre Test'){
                    steps {
                        echo 'Pulling Dependencies'
                        sh 'go version'
                    }
                }

                stage('Test'){

                    steps {
                        echo 'Vetting'

                        sh """cd $GOPATH/src/github.com/cjburchell/reefstatus-go/ && go tool vet ./..."""

                        echo 'Linting'
                        sh """cd $GOPATH/src/github.com/cjburchell/reefstatus-go/ && golint ./..."""

                        echo 'Testing'
                        sh """cd $GOPATH/src/github.com/cjburchell/reefstatus-go/ && go test -race -cover ./..."""
                    }
                }

                stage('Build') {
                    steps {
                        echo 'Building Executable'

                        sh """cd $GOPATH/src/github.com/cjburchell/reefstatus-go/ && go build -o service"""
                    }
                }
        }
}