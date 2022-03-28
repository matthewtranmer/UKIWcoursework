function verticalMenu() {
    var element = document.getElementById("nav-top")

    if (element.className == "nav-top") {
        element.className += " burger"
        return
    } 
    
    element.className = "nav-top"
} 