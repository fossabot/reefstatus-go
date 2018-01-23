pipeline {
        agent {
            docker {
                image 'golang'
                args '-p 3000:3000'
            }
        }

        environment {
                GOPATH = '/var/jenkins_home/workspace/ReefStatus'
                PATH =  "${GOPATH}/bin:$PATH"
        }

        stages {
                stage('Pre Test'){
                    steps {
                        echo "${GOPATH}"
                        echo 'Pulling Dependencies'
                        sh 'go version'
                    }
                }

                stage('Test'){

                    steps {
                        echo 'Vetting'

                        //sh "cd ${GOPATH}/src/github.com/cjburchell/reefstatus-go/ && go tool vet ."

                        echo 'Linting'
                        //sh "cd ${GOPATH}/src/github.com/cjburchell/reefstatus-go/ && golint ."

                        echo 'Testing'
                       // sh "cd ${GOPATH}/src/github.com/cjburchell/reefstatus-go/ && go test -race -cover ."
                    }
                }

                stage('Build') {
                    steps {
                        echo 'Building Executable'

                        sh "cd ${GOPATH}/src/github.com/cjburchell/reefstatus-go/ && go build -o service"
                    }
                }

                stage('Build Image') {
                    steps {
                         echo 'Build Image'
                        sh "cd ${GOPATH}/src/github.com/cjburchell/reefstatus-go/ && docker build -t reefstatus:latest ."
                    }
                }

                stage('Push Image') {
                    steps {
                        echo 'Push Image'
                    }
                }
        }
}