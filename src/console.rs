use docopt::Docopt;

use super::result::Result;
use super::{app, env};

#[derive(Debug, Deserialize)]
struct Args {
    flag_version: bool,
    flag_https: bool,
    flag_name: String,

    cmd_generate: bool,
    cmd_locale: bool,
    cmd_migration: bool,
    cmd_config: bool,
    cmd_nginx: bool,

    cmd_database: bool,
    cmd_create: bool,
    cmd_connect: bool,
    cmd_migrate: bool,
    cmd_rollback: bool,
    cmd_drop: bool,
    cmd_status: bool,

    cmd_start: bool,
}

pub fn run() -> Result<()> {
    let usage = format!(
        "
{name} - {description}.

VERSION: {version}
AUTHORS: {authors}
HOMEPAGE: {homepage}

USAGE:
  {name} generate config
  {name} generate (locale|migration) [--name=<fn>]
  {name} generate nginx [--https]
  {name} database (create|connect|migrate|rollback|status|drop)
  {name} start
  {name} (-h | --help)
  {name} --version

OPTIONS:
  -h --help     Show this screen.
  --version     Show version.
  --name=<fn>   File's name.
  --https       Using https?
  --daemon      Run as daemon mode?
    ",
        version = env::VERSION,
        name = env::NAME,
        description = env::DESCRIPTION,
        homepage = env::HOMEPAGE,
        authors = env::AUTHORS,
    );
    let args: Args = try!(try!(Docopt::new(usage)).deserialize());
    // println!("{:?}", args);
    let app = app::App {};

    if args.flag_version {
        return app.show_version();
    }
    if args.cmd_start {
        return app.start();
    }
    if args.cmd_generate {
        if args.cmd_config {
            return app.generate_config();
        }
        if args.cmd_nginx {
            return app.generate_nginx(args.flag_https);
        }
        if args.cmd_migration {
            return app.generate_migration(args.flag_name);
        }
        if args.cmd_locale {
            return app.generate_locale(args.flag_name);
        }
    }
    if args.cmd_database {
        if args.cmd_create {
            return app.database(|c| c.create(), false);
        }
        if args.cmd_connect {
            return app.database_connect();
        }
        if args.cmd_migrate {
            return app.database_migrate();
        }
        if args.cmd_rollback {
            return app.database_rollback();
        }
        if args.cmd_status {
            return app.database_status();
        }
        if args.cmd_drop {
            return app.database(|c| c.drop(), false);
        }
    }

    return Ok(());
}
