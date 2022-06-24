
import { initializeApp } from "https://www.gstatic.com/firebasejs/9.2.0/firebase-app.js"
import { getMessaging,getToken,deleteToken  } from "https://www.gstatic.com/firebasejs/9.2.0/firebase-messaging.js"

const firebaseConfig = {
    apiKey:"AIzaSyAL3kgQp-cI4Z3acI4ZgPjRoqugmUTRrKY",
    authDomain:"my-demo1-cd429.firebaseapp.com",
    projectId:"my-demo1-cd429",
    storageBucket:"my-demo1-cd429.appspot.com",
    messagingSenderId:"2974294779",
    appId:"1:2974294779:web:523cba077a151b1c44c65e",
    measuermentId:"G-6KVGRHHDHP"
};
const app = initializeApp(firebaseConfig);
const messaging =getMessaging(app);

function requestPermissions() {
    console.log('Requesting permission...');
    Notification.requestPermission().then((permission) => {
      if (permission === 'granted') {
        console.log('Notification permission granted.');
        resetUI();
      } else {
        console.log('Unable to get permission to notify.');
      }
    });
  }

  function deleteOldToken() {
    getToken(messaging, {vapidKey: 'BHCHwpXX2RFjT16mrYaP6dd51qKCbxDCmKB-LjBCfq6y0uYvmyLKzvkmEYcu3iJgcBAXe3QuWRpZbAOj5SDetvI'}).then((currentToken) => {
        deleteToken(messaging,currentToken).then(()=>{
            resetUI();
        }).catch((err)=>{
            console.error("delete token error, ", err)
        })
      }).catch((err) => {
        console.log('Error retrieving registration token. ', err);
        showToken('Error retrieving registration token. ', err);
      });
    }
  
    function resetUI() {
        clearMessages();
        showToken('loading...');
        getToken(messaging, {vapidKey: 'BHCHwpXX2RFjT16mrYaP6dd51qKCbxDCmKB-LjBCfq6y0uYvmyLKzvkmEYcu3iJgcBAXe3QuWRpZbAOj5SDetvI'}).then((currentToken)=>{
            if (currentToken){
                //push token
                var permission_div = document.getElementById("permission_div")
                permission_div.style.display = "none"
                var token_div = document.getElementById("token_div")
                token_div.style.display = "block"
                showToken(currentToken)
            }else{
                console.log("token is empty")
            }
        }).catch((err)=>{
            console.error(err);
        });
      }

      function showToken(currentToken) {
        document.querySelector("#token").textContent = currentToken;
      }
    
      function clearMessages() {
        const messagesElement = document.querySelector('#messages');
        while (messagesElement.hasChildNodes()) {
          messagesElement.removeChild(messagesElement.lastChild);
        }
      }

      function pushToken(currentToken){
        var userId = document.getElementById("userId").textContent
        var param = '{"userId":"'+userId+'","token":"'+currentToken+'","platform":"web"}' 
        console.log("param is ", param);
        const Http = new XMLHttpRequest();
        const url = "/registry"
        Http.open("POST",url);
        Http.send(param);
        Http.onreadystatechange=function(){
          if (this.readyState == 4){
            if ( this.status == 200){
              window.alert("Push token successful!")
            }else{
              window.alert("Push token failure")
            }
          }
        }
      }

      // onMessage(messaging,(payload) =>{
      //   alert(payload)
      //   console.log('Message received. ', payload);
      //   appendMessage(payload);
      //   self.registration.showNotification(notificationTitle,
      //   notificationOptions);
      // });

      // function appendMessage(payload) {
      //   console.log("========append message......")
      //   const messagesElement = document.querySelector('#messages');
      //   const dataHeaderElement = document.createElement('h5');
      //   const dataElement = document.createElement('pre');
      //   dataElement.style = 'overflow-x:hidden;';
      //   dataHeaderElement.textContent = 'Received message:';
      //   dataElement.textContent = JSON.stringify(payload, null, 2);
      //   messagesElement.appendChild(dataHeaderElement);
      //   messagesElement.appendChild(dataElement);
      // }
    

    export{requestPermissions,deleteOldToken,pushToken};