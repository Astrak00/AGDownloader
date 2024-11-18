use std::env;
use std::process::{Child, Command};
use thirtyfour::{DesiredCapabilities, WebDriver};

fn launch_geckodriver(geckodriver_path: &str) -> std::io::Result<Child> {
    // Add a spinner to indicate progress
    let spinner = indicatif::ProgressBar::new_spinner();
    spinner.enable_steady_tick(std::time::Duration::from_millis(100));
    spinner.set_message("Installing geckodriver... (cargo install geckodriver)");

    Command::new("cargo")
        .arg("install")
        .arg("geckodriver")
        .output()
        .expect("failed to install geckodriver");

    Command::new("killall")
        .arg("geckodriver")
        .output()
        .expect("failed to kill geckodriver");

    println!("Launching geckodriver...");
    // Launch geckodriver as a subprocess
    let child = Command::new(geckodriver_path)
        .arg("--port")
        .arg("4444") // Default WebDriver port
        .stdout(std::process::Stdio::null())
        .stderr(std::process::Stdio::null())
        .spawn();

    match child {
        Err(e) => {
            println!("Error launching geckodriver: {:?}", e);
            println!("Please ensure that geckodriver is installed and available in your PATH");
            println!("You may download geckodriver from: https://github.com/mozilla/geckodriver/releases or install it by running cargo install gecodriver");
            Err(e)
        }
        Ok(child) => {
            println!("geckodriver launched successfully");
            Ok(child)
        }
    }
}

pub async fn get_auth_cookie() -> Result<String, Box<dyn std::error::Error>> {
    let geckodriver_path = env::var("GECKODRIVER_PATH").unwrap_or("geckodriver".to_string());
    let mut _geckodriver = launch_geckodriver(&geckodriver_path);
    if _geckodriver.is_err() {
        return Err(Box::new(std::io::Error::new(
            std::io::ErrorKind::Other,
            "Error launching geckodriver",
        )));
    }

    // Wait for geckodriver to start
    std::thread::sleep(std::time::Duration::from_secs(1));

    // Set up Chrome options for a non-headless browser
    let caps = DesiredCapabilities::firefox();

    // Start the WebDriver session
    let driver = WebDriver::new("http://localhost:4444", caps).await?;
    // Add a cookie to the browser session

    // Navigate to a sample webpage
    driver.get("https://aulaglobal.uc3m.es").await?;

    if driver.windows().await?.is_empty() {
        println!("Window closed---------------------");
    }

    loop {
        let cookies = driver.get_all_cookies().await?;
        if cookies
            .iter()
            .any(|cookie| cookie.name == "MoodleSessionag")
        {
            let cookie_moodle = cookies
                .iter()
                .find(|cookie| cookie.name == "MoodleSessionag")
                .unwrap();

            // Close the WebDriver session
            driver.quit().await?;
            _geckodriver.unwrap().kill().unwrap();

            return Ok(cookie_moodle.value.clone());
        }
        tokio::time::sleep(std::time::Duration::from_millis(750)).await;
    }
}
