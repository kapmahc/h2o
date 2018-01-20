use docopt::Docopt;

use super::result::Result;

#[derive(Debug, Deserialize)]
struct Args {
    flag_version: bool,
    flag_https: bool,
    flag_name: String,
    flag_daemon: bool,

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
    cmd_stop: bool,
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
  {name} start [--daemon]
  {name} stop
  {name} (-h | --help)
  {name} --version

OPTIONS:
  -h --help     Show this screen.
  --version     Show version.
  --speed=<kn>  Speed in knots [default: 10].
  --moored      Moored (anchored) mine.
  --drifting    Drifting mine.
    ",
        version = env!("CARGO_PKG_VERSION"),
        name = env!("CARGO_PKG_NAME"),
        description = env!("CARGO_PKG_DESCRIPTION"),
        homepage = env!("CARGO_PKG_HOMEPAGE"),
        authors = env!("CARGO_PKG_AUTHORS")
    );
    let args: Args = try!(try!(Docopt::new(usage)).deserialize());
    if args.flag_version {
        println!("{}", version = env!("CARGO_PKG_VERSION"));
        return Ok(());
    }

    println!("{:?}", args);
    return Ok(());
}
