
function $(id){
    return typeof id === 'string' ? document.getElementById(id):id;
}
window.onload = function(){
    var titles = $('tab-header').getElementsByTagName('li');
    var divs = $('tab-content').getElementsByClassName('dom');
    if(titles.length != divs.length) return;
    for(var i=0; i<titles.length; i++){
        var li = titles[i];
        li.id = i;
        li.onmousemove = function(){
            for(var j=0; j<titles.length; j++){
                titles[j].className = '';
                divs[j].style.display = 'none';
            }
            this.className = 'selected';
            divs[this.id].style.display = 'block';
        }
    }
 
}