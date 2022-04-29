
window.addEventListener('resize', splash);
splash()

function splash(e){
    var width = (window.innerWidth > 0) ? window.innerWidth : screen.width;
    document.getElementById("splashtext").setAttribute("style", "font-size: " + width/2.7 +"%;")
    
    
    if(width > 1400){
        document.getElementById("textbox").setAttribute("style", "display: inline-block; width: 45%;")
        document.getElementById("map").setAttribute("style", "display: inline-block; width: 47%; float: right; margin-right: 30px;")
        document.getElementById("canalimage").setAttribute("style", "float: right; display: inline; width: 45%; padding-right: 31px;")
    }
    else{
        document.getElementById("textbox").setAttribute("style", "display: block; width: 85%; margin: auto; margin-bottom: 20px")
        document.getElementById("map").setAttribute("style", "display: block; width: 87%; float: none; margin: auto;")
        document.getElementById("canalimage").setAttribute("style", "float: none; display: block; width: 85%; padding-right: 0px;")
    }

    if (width > 1000){
        document.getElementById("splashimg").setAttribute("src", "/static/img/homesplash.jpg")
    }
    else{
        document.getElementById("splashimg").setAttribute("src", "/static/img/smallersplash.jpg")
    }
    
}
