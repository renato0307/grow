pub struct DisplayDriver {
    width: u16,
    height: u16,
    rotation: super::Rotation,
}

impl DisplayDriver {
    pub fn new(width: u16, height: u16, rotation: super::Rotation) -> Self {
        DisplayDriver {
            width,
            height,
            rotation,
        }
    }
}

/// The ST7789 driver supports both Parallel and Serial (SPI) ST7789 displays and is intended for use with:
/// - Pico Display
/// - Pico Display 2.0
/// - Tufty 2040
/// - Pico Explorer
/// - 240x240 Round & Square SPI LCD Breakouts
///
/// Implementation based on https://github.com/pimoroni/pimoroni-pico/tree/main/drivers/st7789.
pub struct ST7786 {
    display_driver: DisplayDriver,

    round: bool,

    cs: u8,
    dc: u8,
    wr_sck: u8,
    rd_sck: u8,
    d0: u8,
    bl: u8,
    vsync: u8,
    parallel_sm: u8,

    parallel_offset: u8,
    st_dma: u8,
}
