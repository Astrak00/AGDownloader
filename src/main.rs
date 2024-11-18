mod cookies;
mod courses;
mod resource;

use clap::Parser;

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

    args.cookie = Some("super_token".to_string());

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


    println!("Cookie: {}", auth_cookie);

    // Launch geckodriver as a subprocess
    //let cookie = Cookie::new("MoodleSessionag", "superstingporaqui");
    //driver.add_cookie(cookie).await?;

    // Obtain all the courses that the user is enrolled in
    let courses = courses::get_courses(auth_cookie).await?;
    println!("Courses: {:?}\n\n", courses.len());

    // Parse the individual contents for each course to download
    let mut count = 0;
    let mut all_resources: Vec<resource::Resource> = Vec::new();
    for mut course in courses {
        count += 1;
        if count < 4 {
            continue;
        }
        if count == 5 {
            break;
        }
        let resources: Vec<resource::Resource>;
        (resources, course.fullname) = resource::get_course_contents(auth_cookie, &mut course).await?;
        std::fs::create_dir_all(format!("{}/{}", save_dir, course.fullname)).unwrap();
        if course.id == "178544" {
            for resource in resources.iter() {
                println!("Resource: {:?}", resource.filename);
                println!("link: {:?}", resource.id);
            }
        }
        all_resources.extend(resources);

    }
    let mut count = 0;
    println!("All resources: {:?}", all_resources.len());
    for resource in all_resources {
        count += 1;
        if count < 4 {
            continue;
        }
        if resource.filename == "Presentacion Parte A teoría File" {
            print!("Skipping");
        }
        resource::download_file_with_original_name(&resource, &save_dir, &auth_cookie).await?;
    }




    Ok(())
}
