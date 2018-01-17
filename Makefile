
build:
	cargo build --release
	strip target/release/h2o

clean:
	cargo clean
