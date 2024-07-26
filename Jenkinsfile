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
                    sh '/home/entrust/Public/docker-compose/docker-compose up -d --build'
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
                    /home/entrust/Public/docker-compose/docker-compose down
                    /home/entrust/Public/docker-compose/docker-compose up -d
                    '''
                }
            }
        }
    }

    post {
        always {
            sh '/home/entrust/Public/docker-compose/docker-compose down'
            cleanWs()
        }
    }
}
