mod cookies;
mod courses;
mod resource;

use clap::Parser;
use futures::future::join_all;
use std::sync::{Arc, Mutex};

#[derive(Default, Parser, Debug)]
#[command(
    author = "Astrak00",
    version,
    about = "Program to download AulaGlobal Contents"
)]
struct Args {
    #[arg(short, long)]
    cookie: Option<String>,
    #[arg(short, long)]
    save_dir: Option<String>,
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut args: Args = Args::parse();
    let auth_cookie: String;

    args.cookie = Some("chgtcl0gjoh6a1s6m3kop59fin".to_string());

    if args.cookie.is_none() {
        let result = cookies::get_auth_cookie().await;
        match result {
            Ok(cookie) => {
                auth_cookie = cookie;
            }
            Err(e) => {
                println!("Error getting cookie: {:?}", e);
                return Ok(());
            }
        }
    } else {
        auth_cookie = args.cookie.unwrap();
    }

    let auth_cookie = auth_cookie.as_str();
    let save_dir = args.save_dir.unwrap_or("courses".to_string());


    // Obtain all the courses that the user is enrolled in
    let courses = courses::get_courses(auth_cookie).await?;
    println!("Courses: {:?}\n\n", courses.len());

    // Parse the individual contents for each course to download
    for mut course in courses {
        let resources: Vec<resource::Resource>;
        (resources, course.fullname) = resource::get_course_contents(auth_cookie, &mut course).await?;
        std::fs::create_dir_all(format!("{}/{}", save_dir, course.fullname)).unwrap();
        for resource in resources {
            resource::download_file_with_original_name(&resource, &save_dir, auth_cookie).await?;
        }
    }



    Ok(())
}
