pipeline {
    agent any

    stages {
        stage('Checkout') {
            steps {
                git 'https://github.com/wazadio/realtime-weather.git'
            }
        }

        stage('Build and Run') {
            steps {
                script {
                    sh 'docker-compose up -d --build'
                }
            }
        }

        // stage('Test') {
        //     steps {
        //         script {
        //             sh 'go test'
        //         }
        //     }
        // }

        stage('Deploy') {
            steps {
                script {
                    sh '''
                    docker-compose down
                    docker-compose up -d
                    '''
                }
            }
        }
    }

    post {
        always {
            sh 'docker-compose down'
            cleanWs()
        }
    }
}
