use crate::courses::Course;
use scraper::{Html, Selector};
use std::fs::File;
use std::io::copy;
use std::path::Path;

pub(crate) struct Resource {
    pub(crate) id: String,
    pub(crate) course_name: String,
}

const MOODLE_COURSE_URL: &str = "https://aulaglobal.uc3m.es/course/view.php?id=";
const MOODLE_CONTENT_URL: &str = "https://aulaglobal.uc3m.es/mod/resource/view.php?id=";

pub async fn get_course_contents(
    auth_cookie: &str,
    course: &mut Course,
) -> Result<(Vec<Resource>, String), Box<dyn std::error::Error + Send + Sync>> {
    let client = reqwest::Client::new();
    let res = client
        .get(format!("{}{}", MOODLE_COURSE_URL, course.id))
        .header("Cookie", format!("MoodleSessionag={}", auth_cookie))
        .send()
        .await?
        .text()
        .await?;

    let document = Html::parse_document(&res);

    // Get the h1 element with the course name
    let selector = Selector::parse("h1").unwrap();
    let course_name = document
        .select(&selector)
        .next()
        .unwrap()
        .text()
        .collect::<Vec<_>>()
        .join(" ");
    let clean_course_name = course_name.replace("/", "-");

    let selector = Selector::parse("a").unwrap();

    let mut contents: Vec<Resource> = Vec::new();
    for element in document.select(&selector) {
        if let Some(link) = element.value().attr("href") {
            if link.starts_with("https://aulaglobal.uc3m.es/mod/resource/") {
                contents.push(Resource {
                    id: link.split("=").last().unwrap().to_string(),
                    course_name: course_name.clone(),
                });
            }
        }
    }

    Ok((contents, clean_course_name))
}

pub async fn download_file_with_original_name(
    resource: &Resource,
    output_dir: &str,
    cookie: &str,
) -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
    let id = &resource.id;
    let output_dir = format!("{}/{}", output_dir, resource.course_name.replace("/", "-"));
    let client = reqwest::Client::new();
    let url = format!("{}{}", MOODLE_CONTENT_URL, id);
    let res = client
        .get(&url)
        .header("Cookie", format!("MoodleSessionag={}", cookie))
        .send()
        .await?;

    if !res.status().is_success() {
        return Err(format!("Failed to download file: {}", res.status()).into());
    }

    // From the response headers, extract the filename, which is in the Content-Disposition header as: inline; filename="filename"
    let content_disposition = res.headers().get("Content-Disposition");

    let result: String = match content_disposition {
        Some(header_value) => match header_value.to_str() {
            Ok(result_str) => result_str.to_owned(), // Convert &str to owned String
            Err(_) => {
                // Handle non-UTF8 characters and return owned String
                let bytes = header_value.as_bytes();
                String::from_utf8_lossy(bytes).into_owned()
            }
        },
        None => String::from("Content-Disposition header not found"),
    };

    let filename =
        percent_encoding::percent_decode_str(&parse_filename_from_content_disposition(result))
            .decode_utf8_lossy()
            .to_string();

    let output_path = Path::new(&output_dir).join(&filename);
    let output_path_str = output_path.to_str().unwrap();

    std::fs::create_dir_all(&output_dir)?;
    let mut file = File::create(&output_path)?;

    // Copy the response body to the file
    let content = res.bytes().await?;
    copy(&mut content.as_ref(), &mut file)?;

    println!("File downloaded to: {}", output_path_str);
    Ok(())
}

// Extract the filename from the Content-Disposition header
fn parse_filename_from_content_disposition(header: String) -> String {
    // Look for `filename=` in the header
    header
        .split("filename=")
        .last()
        .and_then(|part| part.split(';').next()) // Get filename part
        .map(|name| name.trim_matches(['"', '\''].as_ref())) // Remove quotes if any
        .unwrap_or("unknown_file") // Default fallback
        .to_string()
}
