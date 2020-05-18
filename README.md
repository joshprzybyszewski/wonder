# wonder 

wonder is intended to increase your desire to see the world around you.

### What is there to wonder about?
Your brain is smart. It takes shortcuts. For example, you likely don't need to think about how to start your car and drive to work: your brain has trained you to do that without hardly thinking. Another shortcut your brain takes is with color. Science tells us each [color is a different wavelength](https://en.wikipedia.org/wiki/Color). White light is different: it contains all of the wavelengths of the visible spectrum. However, when we look at white light, our brains don't say "I see red, orange, yellow, green, blue, indigo, and purple" but rather say "I see _white_". It's a shortcut that our brains take.

### How does this matter in my day to day life?
One way to interpret this shortcut is to realize that our lives and the world around us have so much more information than we are able to interpret or perceive. Another thing to keep in mind is that screens may do a good job of mimicking the real world, but most screens only have red/green/blue pixels to show you. That means they operate because of the shortcuts your brain takes. It also means the next time you see a yellow banana, that's scientifically impossible for your phone to capture and relay that image back to you (at the same wavelength).

To me, this also was enough inspiration to dig into golang + WebAssembly and build a weekend project to filter out the colors that your screen displays back to you. This project is an attempt to wonder at the world around us and to learn a little more code.

### How to run

1. Clone this repo.
1. Run `make build-prod serve`. This builds the wasm output and then serves the minimal html locally.
1. Go to [http://localhost:3434/](http://localhost:3434/).
1. This will ask for access to your webcam. Now you can change the radio buttons to modify the color of the video.
   - Additionally, there's currently an option to "upload an image" which will apply the effect once to that image.

#### Many thanks to the internet

This work started as a copy of [shimmer](github.com/agnivade/shimmer), and then I mixed it with the tutorial [here](https://developer.mozilla.org/en-US/docs/Web/Guide/Audio_and_video_manipulation) to learn about webcam interaction.
