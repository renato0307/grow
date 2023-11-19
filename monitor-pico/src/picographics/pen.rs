extern crate alloc;
use alloc::boxed::Box;

#[allow(dead_code)]
pub enum PicoGraphicsPenType {
    Pen1bit = 0,
    Pen3bit,
    PenP2,
    PenP4,
    PenP8,
    PenRgb332,
    PenRgb565,
    PenRgb888,
    PenInky7,
}

pub fn from(pen_type: PicoGraphicsPenType) -> Result<Box<dyn PicoGraphicsPen>, super::Error> {
    match pen_type {
        PicoGraphicsPenType::PenRgb565 => Ok(Box::new(PicoGraphicsPenRgb565 {})),
        _ => Err(super::Error::NotImplemented),
    }
}

pub trait PicoGraphicsPen {
    fn create_pen(&self, r: u8, g: u8, b: u8) -> u16;
}

struct PicoGraphicsPenRgb565 {}

impl PicoGraphicsPen for PicoGraphicsPenRgb565 {
    fn create_pen(&self, r: u8, g: u8, b: u8) -> u16 {
        (((r & 0b11111000) << 8) | ((g & 0b11111100) << 3) | ((b & 0b11111000) >> 3)).into()
    }
}
