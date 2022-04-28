
window.addEventListener('resize', splash);
splash()

function splash(e){
    var width = (window.innerWidth > 0) ? window.innerWidth : screen.width;
    document.getElementById("splashtext").setAttribute("style", "font-size: " + width/2.7 +"%;")
    
    
    if(width > 1400){
        
        document.getElementById("textbox").setAttribute("style", "display: inline-block; width: 45%;")
        document.getElementById("map").setAttribute("style", "display: inline-block; width: 47%; float: right; margin-right: 30px;")
    }
    else{
        
        document.getElementById("textbox").setAttribute("style", "display: block; width: 85%; margin: auto; margin-bottom: 20px")
        document.getElementById("map").setAttribute("style", "display: block; width: 87%; float: none; margin: auto;")
    }

    if (width > 1000){
        document.getElementById("splashimg").setAttribute("src", "/static/img/homesplash.jpg")
    }
    else{
        document.getElementById("splashimg").setAttribute("src", "/static/img/smallersplash.jpg")
    }
    
}
