//! # UART Example
//!
//! This application demonstrates how to use the UART Driver to talk to a serial
//! connection.
//!
//! It may need to be adapted to your particular board layout and/or pin
//! assignment.
//!
//! See the `Cargo.toml` file for Copyright and license details.

#![no_std]
#![no_main]

// Ensure we halt the program on panic (if we don't mention this crate it won't
// be linked)
use panic_halt as _;

// GPIO traits
use embedded_hal::digital::v2::OutputPin;

// Ensure we halt the program on panic (if we don't mention this crate it won't
// be linked)
use panic_halt as _;

// A shorter alias for the Peripheral Access Crate, which provides low-level
// register access
use pimoroni_pico_lipo_16mb::hal::pac;

// A shorter alias for the Hardware Abstraction Layer, which provides
// higher-level drivers.
use pimoroni_pico_lipo_16mb::hal;

// Alias for our HAL crate
use pimoroni_pico_lipo_16mb as bsp;

// A shorter alias for the Peripheral Access Crate, which provides low-level
// register access
use bsp::entry;

// Some traits we need
use bsp::hal::clocks::Clock;

/// External high-speed crystal on the Raspberry Pi Pico board is 12 MHz. Adjust
/// if your board has a different frequency
const XTAL_FREQ_HZ: u32 = 12_000_000u32;

mod picographics;

/// Allocator
extern crate alloc;

use embedded_alloc::Heap;

#[global_allocator]
static HEAP: Heap = Heap::empty();

/// Entry point to our bare-metal application.
///
/// The `#[rp2040_hal::entry]` macro ensures the Cortex-M start-up code calls this function
/// as soon as all global variables and the spinlock are initialised.
///
/// The function configures the RP2040 peripherals, then writes to the UART in
/// an infinite loop.
#[entry]
fn main() -> ! {
    // Initialize the allocator BEFORE you use it
    {
        use core::mem::MaybeUninit;
        const HEAP_SIZE: usize = 1024;
        static mut HEAP_MEM: [MaybeUninit<u8>; HEAP_SIZE] = [MaybeUninit::uninit(); HEAP_SIZE];
        unsafe { HEAP.init(HEAP_MEM.as_ptr() as usize, HEAP_SIZE) }
    }

    // Grab our singleton objects
    let mut pac = pac::Peripherals::take().unwrap();
    let core = pac::CorePeripherals::take().unwrap();

    // Set up the watchdog driver - needed by the clock setup code
    let mut watchdog = bsp::hal::Watchdog::new(pac.WATCHDOG);

    // Configure the clocks
    let clocks = bsp::hal::clocks::init_clocks_and_plls(
        XTAL_FREQ_HZ,
        pac.XOSC,
        pac.CLOCKS,
        pac.PLL_SYS,
        pac.PLL_USB,
        &mut pac.RESETS,
        &mut watchdog,
    )
    .ok()
    .unwrap();

    let mut delay = cortex_m::delay::Delay::new(core.SYST, clocks.system_clock.freq().to_Hz());

    // Inits our display
    let mut display = picographics::PicoGraphics::new(
        picographics::PicoGraphicsDisplay::DisplayPicoDisplay,
        Some(0),
        None,
        None,
    );

    // The single-cycle I/O block controls our GPIO pins
    let sio = hal::Sio::new(pac.SIO);

    // Set the pins up according to their function on this particular board
    let pins = pimoroni_pico_lipo_16mb::Pins::new(
        pac.IO_BANK0,
        pac.PADS_BANK0,
        sio.gpio_bank0,
        &mut pac.RESETS,
    );

    // Set the LED to be an output
    let mut led_pin = pins.led.into_push_pull_output();
    loop {
        display.set_backlight(0.8).ok().unwrap();
        let white = display.create_pen(255, 255, 255);
        let black = display.create_pen(0, 0, 0);

        display.set_pen(white);
        display.clear();
        delay.delay_ms(1000);

        display.set_pen(black);
        display.clear();
        delay.delay_ms(1000);

        led_pin.set_high().unwrap();
        delay.delay_ms(500);
        led_pin.set_low().unwrap();
        delay.delay_ms(500);
    }
}

// End of file
