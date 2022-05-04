function verticalMenu() {
    var element = document.getElementById("nav-top")

    if (element.className == "nav-top") {
        element.className += " burger"
        return
    } 
    
    element.className = "nav-top"
} 

window.addEventListener('resize', resize);
resize()
function resize(e){
    var width = (window.innerWidth > 0) ? window.innerWidth : screen.width;
    
    console.log(width);
    if(width > 900){
        document.getElementById("mobile-content").setAttribute("style", "display: none;")
        document.getElementById("desktop-content").setAttribute("style", "display: block;")
    }
    else{
        document.getElementById("desktop-content").setAttribute("style", "display: none;")
        document.getElementById("mobile-content").setAttribute("style", "display: block;")
    }
}
