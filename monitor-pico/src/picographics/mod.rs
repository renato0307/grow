pub mod pen;
pub mod st7786;

extern crate alloc;

use alloc::boxed::Box;
use core::option::Option;
use core::option::Option::Some;
use core::result::Result;
use core::result::Result::Err;
use core::result::Result::Ok;

use pen::*;
use st7786::*;

pub enum Error {
    /// Feature not implemented
    NotImplemented,

    /// Brightness out of range - expected 0.0 to 1.0
    BrightnessOutOfRange,
}

#[allow(dead_code)]
pub enum Rotation {
    Rotate0 = 0,
    Rotate90 = 90,
    Rotate180 = 180,
    Rotate270 = 270,
}

#[allow(dead_code)]
pub enum PicoGraphicsDisplay {
    DisplayPicoDisplay,
}

#[allow(dead_code)]
pub enum PicoGraphicsBusType {
    BusSpi,
}

impl Into<i8> for Rotation {
    fn into(self) -> i8 {
        self as i8
    }
}

pub struct PicoGraphics {
    bus: PicoGraphicsBusType,
    display: PicoGraphicsDisplay,
    pen_type: Box<dyn pen::PicoGraphicsPen>,
    rotate: i8,

    color: u16,
}

impl PicoGraphics {
    // The constructor of a PicoGraphics
    pub fn new(
        display: PicoGraphicsDisplay,
        rotate: Option<i8>,
        bus: Option<PicoGraphicsBusType>,
        pen_type: Option<pen::PicoGraphicsPenType>,
    ) -> Self {
        let rotate: i8 = match rotate {
            Some(v) => v,
            None => -1,
        };
        let pt = match pen_type {
            Some(pt) => pt,
            None => PicoGraphicsPenType::PenRgb565,
        };
        let result = pen::from(pt).ok().unwrap();

        let mut g = PicoGraphics {
            bus: match bus {
                Some(b) => b,
                None => PicoGraphicsBusType::BusSpi,
            },
            display,
            pen_type: result,
            rotate,
            color: 0,
        };

        let (valid, width, height, mut rotate, bus_type) = g.get_display_settings();
        if !valid {
            panic!("Invalid display")
        }

        if rotate == -1 {
            rotate = Rotation::Rotate0.into();
        }

        g.rotate = rotate;
        g.bus = bus_type;
        g.set_pen(0);
        g.clear();

        return g;
    }

    pub fn set_backlight(&self, brightness: f32) -> Result<(), Error> {
        if brightness < 0.0 || brightness > 1.0 {
            return Err(Error::BrightnessOutOfRange);
        }
        return Ok(());
    }

    pub fn create_pen(&self, r: u8, g: u8, b: u8) -> u16 {
        self.pen_type.create_pen(r, g, b)
    }

    pub fn set_pen(&mut self, pen: u16) {
        self.color = pen
    }

    fn get_display_settings(
        &self,
    ) -> (
        bool,                // valid
        u8,                  // width
        u8,                  // height
        i8,                  // rotate
        PicoGraphicsBusType, // bus_type
    ) {
        match self.display {
            PicoGraphicsDisplay::DisplayPicoDisplay => (
                true,
                240,
                135,
                match self.rotate {
                    -1 => Rotation::Rotate270.into(),
                    _ => self.rotate as i8,
                },
                PicoGraphicsBusType::BusSpi,
            ),
            _ => (false, 0, 0, 0, PicoGraphicsBusType::BusSpi),
        }
    }

    // clears the screen by setting the pen to black and drawing a rectangle
    pub fn clear(&self) {}

    pub fn rectangle(&self, x0: i32, y0: i32, x1: i32, y1: i32) {}
}
