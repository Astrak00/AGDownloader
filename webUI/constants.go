package webui

const (
	correctResponseHTML = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Selected Courses</title>
    <style>
      * {
        box-sizing: border-box;
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
          Oxygen, Ubuntu, Cantarell, "Open Sans", "Helvetica Neue", sans-serif;
      }

      body {
        max-width: 800px;
        margin: 0 auto;
        padding: 20px;
        background-color: #f5f5f5;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        min-height: 100vh;
        text-align: center;
      }

      .confirmation-card {
        background-color: white;
        border-radius: 12px;
        box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
        padding: 40px;
        width: 100%;
        max-width: 500px;
        margin-bottom: 30px;
        border-top: 5px solid #2563eb;
      }

      h1 {
        color: #333;
        margin-top: 0;
        margin-bottom: 20px;
        font-size: 28px;
      }

      .icon-container {
        margin: 20px 0;
      }

      .checkmark-circle {
        width: 80px;
        height: 80px;
        border-radius: 50%;
        background-color: #eff6ff;
        display: flex;
        align-items: center;
        justify-content: center;
        margin: 0 auto;
      }

      .checkmark {
        color: #2563eb;
        font-size: 40px;
      }

      p {
        color: #666;
        font-size: 16px;
        line-height: 1.6;
        margin-bottom: 30px;
      }

      .message {
        margin-bottom: 30px;
      }

      .close-button {
        display: inline-block;
        padding: 12px 24px;
        background-color: #2563eb;
        color: white;
        border: none;
        border-radius: 6px;
        font-size: 16px;
        font-weight: 600;
        cursor: pointer;
        transition: background-color 0.2s ease;
        text-decoration: none;
      }

      .close-button:hover {
        background-color: #1d4ed8;
      }

      .countdown {
        font-size: 14px;
        color: #888;
        margin-top: 20px;
      }

      @media (max-width: 600px) {
        .confirmation-card {
          padding: 30px 20px;
        }
      }
    </style>
    <script type="text/javascript">
      window.onload = function () {
        // Attempt to close the current tab
        try {
          window.open("", "_self", "");
          window.close();
        } catch (e) {
          console.log("Could not automatically close the window:", e);
        }

        // Countdown timer as fallback
        let seconds = 5;
        const countdownElement = document.getElementById("countdown");

        const interval = setInterval(function () {
          seconds--;
          countdownElement.textContent = seconds;

          if (seconds <= 0) {
            clearInterval(interval);
            countdownElement.parentElement.textContent =
              "Please close this tab manually.";
          }
        }, 1000);
      };

      function closeTab() {
        window.close();
      }
    </script>
  </head>
  <body>
    <div class="confirmation-card">
      <h1>Courses Selected Successfully</h1>

      <div class="icon-container">
        <div class="checkmark-circle">
          <span class="checkmark">✓</span>
        </div>
      </div>

      <div class="message">
        <p>
          Your course selection has been submitted successfully. Please close
          this tab.
        </p>
      </div>

      <div class="countdown">
        Closing in <span id="countdown">5</span> seconds...
      </div>
    </div>
  </body>
</html>
`

	courseSelectorHTMLStart = `<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Course Selector</title>
    <style>
        * {
            box-sizing: border-box;
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
        }
        
        body {
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        
        h1 {
            text-align: center;
            color: #333;
            margin-bottom: 30px;
        }
        
        .course-container {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
            gap: 20px;
            margin-bottom: 30px;
        }
        
        .course-option {
            position: relative;
            height: 150px;
            border: 2px solid #e0e0e0;
            border-radius: 8px;
            background-color: white;
            padding: 20px;
            cursor: pointer;
            transition: all 0.2s ease;
            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;
            text-align: center;
            box-shadow: 0 2px 4px rgba(0,0,0,0.05);
        }
        
        .course-option:hover {
            border-color: #2563eb;
            transform: translateY(-2px);
            box-shadow: 0 4px 8px rgba(0,0,0,0.1);
        }
        
        .course-option.selected {
            border-color: #2563eb;
            background-color: #eff6ff;
        }
        
        .course-title {
            font-weight: 600;
            margin-bottom: 8px;
            color: #333;
        }
        
        .course-description {
            font-size: 14px;
            color: #666;
        }
        
        .checkmark {
            position: absolute;
            top: 10px;
            right: 10px;
            width: 20px;
            height: 20px;
            border-radius: 50%;
            background-color: #2563eb;
            color: white;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 12px;
            opacity: 0;
            transition: opacity 0.2s ease;
        }
        
        .course-option.selected .checkmark {
            opacity: 1;
        }
        
        .hidden-checkbox {
            position: absolute;
            opacity: 0;
            cursor: pointer;
            height: 0;
            width: 0;
        }
        
        .submit-button {
            display: block;
            width: 100%;
            max-width: 200px;
            margin: 0 auto;
            padding: 12px 20px;
            background-color: #2563eb;
            color: white;
            border: none;
            border-radius: 6px;
            font-size: 16px;
            font-weight: 600;
            cursor: pointer;
            transition: background-color 0.2s ease;
			margin-bottom: 30px;
        }
        
        .submit-button:hover {
            background-color: #1d4ed8;
        }
        
        @media (max-width: 600px) {
            .course-container {
                grid-template-columns: 1fr;
            }
        }
    </style>
</head>
<body>
    <h1>Course Selector</h1>

	<button type="button" id="select-all" class="submit-button">Select All</button>
    
    <form action="/submit" method="post">
        <div class="course-container">`

	courseSelectorHTMLEnd = `</div>
        <button type="submit" value="Submit" class="submit-button">Submit</button>
    </form>

     <script>
        document.addEventListener('DOMContentLoaded', function() {
            const courseOptions = document.querySelectorAll('.course-option');
            const selectAllButton = document.getElementById('select-all');

            courseOptions.forEach(option => {
                option.addEventListener('click', function(event) {
                    const checkbox = this.querySelector('input[type="checkbox"]');
                    checkbox.checked = !checkbox.checked;
                    this.classList.toggle('selected', checkbox.checked);
                    event.preventDefault();
                });
            });

            selectAllButton.addEventListener('click', function() {
                let allSelected = [...courseOptions].every(option => option.querySelector('input').checked);
                
                courseOptions.forEach(option => {
                    const checkbox = option.querySelector('input[type="checkbox"]');
                    checkbox.checked = !allSelected;
                    option.classList.toggle('selected', !allSelected);
                });

                selectAllButton.textContent = allSelected ? "Select All" : "Deselect All";
            });
        });
    </script>
</body>
</html>`
)
