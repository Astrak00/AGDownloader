use scraper::{Html, Selector};

#[derive(Clone)]
pub(crate) struct Course {
    pub(crate) id: String,
    pub(crate) fullname: String,
}

const MOODLE_URL: &str = "https://aulaglobal.uc3m.es/";

pub async fn get_courses(auth_cookie: String) -> Result<Vec<Course>, Box<dyn std::error::Error>> {
    let client = reqwest::Client::new();
    let res = client
        .get(MOODLE_URL)
        .header("Cookie", format!("MoodleSessionag={}", auth_cookie))
        .send()
        .await?
        .text()
        .await?;

    let document = Html::parse_document(&res);

    let selector = Selector::parse("p.coursename > a").unwrap();

    let mut courses: Vec<Course> = Vec::new();
    for element in document.select(&selector) {
        if let Some(link) = element.value().attr("href") {
            courses.push(Course {
                id: link.split("=").last().unwrap().to_string(),
                fullname: "".to_string(),
            });
        }
    }

    Ok(courses)
}
