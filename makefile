node: src/**.hx
	@haxe build/node.hxml
	@pkg -t node12-linux --out-path dist/ .

java: src/**.hx
	@haxe build/java.hxml

clean: bin/** dist/**
	@rm -rf bin dist
