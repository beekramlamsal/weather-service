fetch('api')
  .then(res => res.text())
  .then(text => {
    const display = document.getElementById('weather-text');
    const icon = document.getElementById('weather-icon');
    display.innerText = text;

    const desc = text.toLowerCase();

    if (desc.includes("sun") || desc.includes("clear")) {
      icon.src = "images/sun.png";
    } else if (desc.includes("cloud")) {
      icon.src = "images/cloud.png";
    } else if (desc.includes("rain")) {
      icon.src = "images/rain.png";
    } else if (desc.includes("thunder") || desc.includes("storm")) {
      icon.src = "images/thunder.png";
    } else if (desc.includes("wind")) {
      icon.src = "images/wind.png";
    } else {
      icon.src = "images/default.png";
    }

    icon.hidden = false;
  })
  .catch(err => {
    document.getElementById('weather-text').innerText = "Error fetching weather.";
    console.error(err);
  });
