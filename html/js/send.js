
function sendMessageToUser(userId){
    var param = '{"userId":"'+userId+'"}'
    console.log(param)
    const Http = new XMLHttpRequest();
    const url = "/sendMessageToUser"
    Http.open("POST",url);
    Http.send(param);
    Http.onreadystatechange=function(){
      if (this.readyState == 4){
        if ( this.status == 200){
          window.alert("send message successful!")
        }else{
          window.alert("send message failure!", this.responseText)
        }
      }
    }
}

function sendMessageToTopic(topic){
  var param = '{"topic":"'+topic+'"}'
  console.log(param)
  const Http = new XMLHttpRequest();
  const url = "/sendMessageToTopic"
  Http.open("POST",url);
  Http.send(param);
  Http.onreadystatechange=function(){
    if (this.readyState == 4){
      if ( this.status == 200){
        window.alert("send message successful!")
      }else{
        window.alert("send message failure!", this.responseText)
      }
    }
  }
}

function sendMessageToPlatform(platform){
  var param = '{"platform":"'+platform+'"}'
  console.log(param)
  const Http = new XMLHttpRequest();
  const url = "/sendMessageToPlatform"
  Http.open("POST",url);
  Http.send(param);
  Http.onreadystatechange=function(){
    if (this.readyState == 4){
      if ( this.status == 200){
        window.alert("send message successful!")
      }else{
        window.alert("send message failure!", this.responseText)
      }
    }
  }
}

function subscribeTopic(topic, userId){
  var param='{"userId":["'+userId+'"],"topic":"'+topic+'"}'
  console.log(param)
  const Http = new XMLHttpRequest();
    const url = "/subscribeTopic"
    Http.open("POST",url);
    Http.send(param);
    Http.onreadystatechange=function(){
      if (this.readyState == 4){
        if ( this.status == 200){
          window.alert("subscribe topic successful!")
        }else{
          alert("subscribe topic failure!", this.responseText)
        }
      }
    }
}

function unsubscribeTopic(topic, userId){
  var param='{"userId":["'+userId+'"],"topic":"'+topic+'"}'
  console.log(param)
  const Http = new XMLHttpRequest();
    const url = "/unsubscribeTopic"
    Http.open("POST",url);
    Http.send(param);
    Http.onreadystatechange=function(){
      if (this.readyState == 4){
        if ( this.status == 200){
          window.alert("unsubscribe topic successful!")
        }else{
          alert("unsubscribe topic failure!", this.responseText)
        }
      }
    }
}

export{sendMessageToUser, subscribeTopic, unsubscribeTopic,sendMessageToTopic,sendMessageToPlatform };