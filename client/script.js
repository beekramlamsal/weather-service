fetch('api')
    .then(res => res.text())
    .then(text => {
        const display = document.getElementById('weather-text');
        const icon = document.getElementById('weather-icon');
        
        // remove loading class and add fade-in effect
        display.classList.remove('loading');
        display.innerText = text;
        
        const desc = text.toLowerCase();

        const firstSentence = desc.split(".")[0];
        
        if (firstSentence.includes("sunny") || firstSentence.includes("clear")) {
            icon.src = "images/sun.png";
        } else if (firstSentence.includes("cloud")) {
            icon.src = "images/cloud.png";
        } else if (firstSentence.includes("rain") || firstSentence.includes("showers")) {
            icon.src = "images/rain.png";
        } else if (firstSentence.includes("thunder") || firstSentence.includes("storm")) {
            icon.src = "images/thunder.png";
        } else if (firstSentence.includes("wind")) {
            icon.src = "images/wind.png";
        } else {
            icon.src = "images/default.png";
        }
        
        icon.hidden = false;
        
        // a slight delay for smooth animation
        setTimeout(() => {
            icon.classList.add('show');
        }, 100);
    })
    .catch(err => {
        const display = document.getElementById('weather-text');
        display.classList.remove('loading');
        display.classList.add('error');
        display.innerText = "Error fetching weather.";
        console.error(err);
    });