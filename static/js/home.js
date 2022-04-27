window.addEventListener('resize', splash);
splash()

function splash(e){
    var width = (window.innerWidth > 0) ? window.innerWidth : screen.width;
    document.getElementById("splashtext").setAttribute("style", "font-size: " + width/2.7 +"%;")

    if(width > 900){
        document.getElementById("splashimg").setAttribute("src", "/static/img/homesplash.jpg")
        document.getElementById("textbox").setAttribute("style", "display: inline-block; width: 40%;")
        document.getElementById("map").setAttribute("style", "display: inline-block; width: 45%; float: right; margin-right: 30px;")
    }
    else{
        document.getElementById("splashimg").setAttribute("src", "/static/img/smallersplash.jpg")
        document.getElementById("textbox").setAttribute("style", "display: block; width: 85%; margin: auto; margin-bottom: 20px")
        document.getElementById("map").setAttribute("style", "display: block; width: 90%; float: none; margin: auto;")
    }
}
